c = "2*(111+2)"

idx = 0

def terminal() -> int:
    global idx
    res = ''
    while idx < len(c) and c[idx].isnumeric():
        res += c[idx]
        idx += 1
    if res != '':
        return int(res)
    if c[idx] == '(':
        idx +=1
        res2 = eval()
        if idx < len(c) and c[idx] == ')':
            return res2
    raise ValueError("Could not parse")

def operation() -> int:
    global idx
    res = terminal()
    while idx < len(c) and c[idx] in ['*', '/']:
        
        if c[idx] == '*':
            idx += 1
            res *= terminal()
        elif c[idx] == '/':
            idx += 1
            res /= terminal()
            
    return int(res)

def eval() -> int:
    global idx
    res = operation()
    while idx < len(c) and c[idx] == '+':
        idx +=1
        res += eval()
    return res
print(c)
print(eval())