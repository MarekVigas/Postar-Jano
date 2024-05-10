from datetime import datetime
import pandas as pd
import os
import re
from nordigen import NordigenClient
import psycopg2

db_conn = psycopg2.connect(database=os.environ.get("DB_NAME"),
                        host=os.environ.get("DB_HOST"),
                        user=os.environ.get("DB_USER"),
                        password=os.environ.get("DB_PASSWORD"),
                        port=os.environ.get("DB_PORT"))

client = NordigenClient(
    secret_id=os.environ.get("NORDIGEN_SECRET_ID"),
    secret_key=os.environ.get("NORDIGEN_SECRET_KEY"),
)

client.generate_token()

requisitions = os.environ.get("NORDIGEN_REQUISITIONS").split(",")
matched_transactions = []
unmatched_transactions = []

for requisition_id in requisitions:
    requisition = client.requisition.get_requisition_by_id(requisition_id)

    account = client.account_api(id=requisition["accounts"][0])

    transactions_raw = account.get_transactions()["transactions"]["booked"]

    bank = requisition["institution_id"]

    for transaction in transactions_raw:
        if transaction.get("endToEndId") is None:
            print("Skipping transaction without endToEndId")
            unmatched_transactions.append(transaction)
            continue
        e2eId = transaction["endToEndId"]

        # regex match
        if bank == "PRIMABANK_KOMASK2X":
            p = re.compile("^\/VS0{0,8}(?P<vs>\d{2})/SS0{0,4}(?P<ss>100\d{3})/KS")
            ref = transaction["transactionId"]
        elif bank == "FIO_FIOZSKBA":
            p = re.compile("^\?VS0{0,8}(?P<vs>\d{2})SS0{0,4}(?P<ss>100\d{3})KS")
            ref = transaction["entryReference"]
        else:
            print("Unknown bank")
            exit(1)

        m = p.match(e2eId)

        if m is None:
            print(f"Skipping transaction with e2eId: {e2eId} ref: {ref}")
            unmatched_transactions.append(transaction)
            continue

        vs = m.group("vs")
        ss = m.group("ss")
        if vs not in ["11", "12", "13", '14']:
            print(f"Skipping transaction with not allowed VS {vs} e2eId: {e2eId} ref: {ref}")
            unmatched_transactions.append(transaction)
            continue

        currency = transaction["transactionAmount"]["currency"]

        if currency != "EUR":
            print(f"Skipping transaction with not allowed currency {currency} ref: {ref}")
            unmatched_transactions.append(transaction)
            continue

        float_amount = float(transaction["transactionAmount"]["amount"])
        int_amount = int(float_amount)

        if float_amount != float(int_amount):
            print(f"Skipping transaction with not whole amount {float_amount} ref: {ref}")
            unmatched_transactions.append(transaction)
            continue

        matched_transactions.append(
            {
                "ref": ref,
                "date": transaction["bookingDate"],
                "amount": int_amount,
                "currency": transaction["transactionAmount"]["currency"],
                "vs": m.group("vs"),
                "ss": m.group("ss"),
            }
        )

cursor = db_conn.cursor()

for transaction in matched_transactions:
    cursor.execute("select * from registrations_with_event where payed is null and specific_symbol = %s and payment_reference = %s", (transaction["ss"], transaction["vs"]))
    row = cursor.fetchone()
    if row is None:
        print(f"Skipping transaction VS: {transaction['vs']} SS: {transaction['ss']}. It is probably already payed.")
        continue
    cursor.execute("update registrations set payed = %s where id = %s", (transaction["amount"], row[0]))
    print(f"Updated registration with id: {row[0]} with amount: {transaction['amount']}")

db_conn.commit()
db_conn.close()

unmatched_df = pd.json_normalize(unmatched_transactions)

unmatched_df['transactionAmount.amount'] = unmatched_df['transactionAmount.amount'].astype(float)
unmatched_df = unmatched_df[unmatched_df['transactionAmount.amount'] > 0]

unmatched_df.to_csv(f"unmatched_{datetime.now()}.csv", encoding='utf-8', index=False)
