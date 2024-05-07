import json
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

requisition = client.requisition.get_requisition_by_id(os.environ.get("NORDIGEN_REQUISITION_ID"))

account = client.account_api(id=requisition["accounts"][0])

transactions_raw = account.get_transactions()["transactions"]["booked"]

transactions = []

for transaction in transactions_raw:
    e2eId = transaction["endToEndId"]

    # regex match
    p = re.compile("^\?VS(?P<vs>\d{2})SS(?P<ss>\d{6})KS")
    m = p.match(e2eId)

    if m is None:
        print(f"Skipping transaction with e2eId: {e2eId} ref: {transaction['entryReference']}")
        continue

    vs = m.group("vs")
    ss = m.group("ss")
    if vs not in ["11", "12", "13", '14']:
        print(f"Skipping transaction with not allowed VS {vs} e2eId: {e2eId} ref: {transaction['entryReference']}")
        continue

    currency = transaction["transactionAmount"]["currency"]

    if currency != "EUR":
        print(f"Skipping transaction with not allowed currency {currency} ref: {transaction['entryReference']}")
        continue

    float_amount = float(transaction["transactionAmount"]["amount"])
    int_amount = int(float_amount)

    if float_amount != float(int_amount):
        print(f"Skipping transaction with not whole amount {float_amount} ref: {transaction['entryReference']}")
        continue

    transactions.append(
        {
            "ref": transaction["entryReference"],
            "date": transaction["bookingDate"],
            "amount": int_amount,
            "currency": transaction["transactionAmount"]["currency"],
            "vs": m.group("vs"),
            "ss": m.group("ss"),
        }
    )

print(json.dumps(transactions, indent=2))

cursor = db_conn.cursor()

for transaction in transactions:
    cursor.execute("select * from registrations_with_event where payed is null and specific_symbol = %s and payment_reference = %s", (transaction["ss"], transaction["vs"]))
    row = cursor.fetchone()
    if row is None:
        print(f"Skipping transaction VS: {transaction['vs']} SS: {transaction['ss']}. It is probably already payed.")
        continue
    cursor.execute("update registrations set payed = %s where id = %s", (transaction["amount"], row[0]))
    print(f"Updated registration with id: {row[0]} with amount: {transaction['amount']}")

db_conn.commit()
db_conn.close()
