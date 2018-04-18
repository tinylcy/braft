import glob

files = glob.glob("./*_data")

sum = 0
lineCount = 0 
for file in files:
    with open(file) as f:
        lines = f.readlines()
        lineCount += len(lines)
        for line in lines:
            sum += int(line[:-1])


print("sum", sum)
print("lineCount", lineCount)
print(float(sum) / float(lineCount), 'ms')