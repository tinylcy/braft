with open('./data') as f:
    lines = f.readlines()

count = 100
sum = 0
for line in lines:
    if line[len(line)-3:] == "ms\n":
        # print(float(line[:len(line)-3]))
        sum += float(line[:len(line)-3])
    else:
        # print(1000 * float(line[:len(line)-2]))
        sum += 1000 * float(line[:len(line)-2])

avg = sum / len(lines)
print(avg, 'ms')
print(count / (avg / 1000), "tps")

    