import os

file_size = 1024 * 1024

file_name = "big_file.txt"

with open(file_name, "wb") as file:
    file.seek(file_size - 1)
    file.write(b"\0")

if os.path.exists(file_name):
    print("FILE CREATED")