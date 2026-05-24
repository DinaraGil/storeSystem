import asyncio
from bleak import BleakClient
import psycopg2

conn = psycopg2.connect(
    host="localhost",
    port=5433,
    database="storagedb",
    user="storageuser",
    password="storagepass"
)

cur = conn.cursor()

def insert_event(rfid_id, article, is_in, error, created_at):
    cur.execute(
        """
        INSERT INTO event (rfid_id, article, is_in, error, created_at)
        VALUES (%s, %s, %s, %s, %s)
        """,
        (rfid_id, article, is_in, error, created_at)
    )
    conn.commit()


address = "20:91:48:66:3B:82"
CHAR_UUID = "0000ffe1-0000-1000-8000-00805f9b34fb"

buffer = ""

def notification_handler(sender, data):
    global buffer

    chunk = data.decode(errors="ignore")
    buffer += chunk
    error = ""
    rfid_id = 0
    article = 0
    is_in = True
    created_at = ""

    while "\n" in buffer:
        line, buffer = buffer.split("\n", 1)
        print("Full message:", line)

        print(line.split(", "))
        data_splited = line.split(", ")
        if "ERROR" in data_splited[0]:
            error = data_splited[0].split("=")[1]
            rfid_id = data_splited[1].split("=")[1]
            created_at = data_splited[2].split("=")[1].strip('\r')
            print(rfid_id, is_in, error, created_at)

            insert_event(rfid_id, None, None, error, created_at)
        else:
            is_in = data_splited[0] == "IN"
            rfid_id = data_splited[1].split("=")[1]
            article_temp = data_splited[2].split("=")[1].replace(".", "")
            article = article_temp.replace("Ten", "")
            created_at = data_splited[3].split("=")[1].strip('\r')
            print(rfid_id, article, is_in, error, created_at)

            insert_event(rfid_id, article, is_in, None, created_at)


async def main():
    async with BleakClient(address) as client:
        await client.start_notify(CHAR_UUID, notification_handler)
        print("Connected. Waiting for data...")

        while True:
            await asyncio.sleep(1)

asyncio.run(main())