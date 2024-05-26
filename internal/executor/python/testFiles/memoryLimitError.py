def allocate_memory(size_in_bytes):
    return bytearray(size_in_bytes)

if __name__ == "__main__":
    memory_size = 1024 * 1024 * 200

    try:
        big_array = allocate_memory(memory_size)
        print("SUCCESS")
    except MemoryError:
        print("FAIL")
        input()

    input()