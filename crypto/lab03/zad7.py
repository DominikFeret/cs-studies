import random

def read_input():
    return int(input())

def print_output(result):
    for byte in result:
        if 32 <= byte <= 126:
            print(chr(byte), end='')
        else:  
            print('\\x{:02x}'.format(byte), end='')

k = read_input()
n = read_input()
random.seed(k)

random_bytes = random.randbytes(n)
print(repr(random_bytes)[2:-1], end='')

print_output(random_bytes)