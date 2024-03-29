def xor_strings(s1, s2):
    return ''.join(chr(ord(a) ^ ord(b)) for a, b in zip(s1, s2))

def read_input():
    try:
        return input().encode('ascii').decode('unicode_escape')
    except UnicodeDecodeError:
        return input()

def print_output(result):
    for char in result:
        if 32 <= ord(char) <= 126:
            print(char, end='')
        else:
            print('\\x{:02x}'.format(ord(char)), end='')


s1 = read_input()
s2 = read_input()

result = xor_strings(s1, s2)

print_output(result)