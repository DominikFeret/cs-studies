import random

def xor_bytes(b1, b2):
    return bytes(a ^ b for a, b in zip(b1, b2))

def read_input():
    string = bytes(input(), 'unicode_escape')
    seed = int(input())
    return string, seed

def print_output(result):
    for byte in result:
        if 32 <= byte <= 126:
            print(chr(byte), end='')
        else:
            print('\\x{:02x}'.format(byte), end='')

string, seed = read_input()
n = len(string)
random.seed(seed)

key = random.randbytes(n)

result = xor_bytes(string, key)
print_output(result)