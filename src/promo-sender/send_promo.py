import logging
import os
import time
import requests

from csv import reader

logging.basicConfig(level=logging.INFO)

base_url = "https://leto-api.salezko.sk"
login_url = f"{base_url}/api/sign/in"
generate_promo_url = f"{base_url}/api/promo_codes"


def str2bool(v):
    return v.lower() in ("yes", "true", "t", "1")

username = os.getenv("USERNAME")
password = os.getenv("PASS")
input_file = os.getenv("INPUT_FILE")
send_email = str2bool(os.getenv("SEND_EMAIL"))


def login(user, psswd):
    response = requests.post(login_url, json={
        "username": user,
        "password": psswd,
    })

    if response.status_code != 200:
        logging.error("Status Code", response.status_code)
        logging.error("JSON Response ", response.json())
        raise Exception("Failed to login")
    logging.info("Logged in successfully.")
    return response.json()['token']


if __name__ == "__main__":
    token = login(username, password)

    with open(input_file, 'r') as read_obj:
        csv_reader = reader(read_obj)
        for row in csv_reader:
            time.sleep(1)

            data = {
                "email": row[1],
                "registration_count": int(row[0]),
                "send_email": bool(send_email),
            }
            response = requests.post(generate_promo_url, json=data, headers={"Authorization": f"Bearer {token}"})
            if response.status_code != 200:
                logging.error(f"Failed to send email: {row[1]}")
                logging.error(f"Status Code {response.status_code}")
                logging.error(f"JSON Response {response.json()}")
                raise Exception("Failed to send email")

            logging.info(f"promo code generated {row[1]} {response.json()['promo_code']}")
