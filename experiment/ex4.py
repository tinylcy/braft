import matplotlib
import matplotlib.pyplot as plt
import numpy as np

x1 = np.linspace(6, 30, 25)
x2 = np.linspace(4, 30, 27)

y1 = [157, 158, 159, 158, 159, 160, 163, 160, 155, 150, 149, 144, 139, 127, 126, 111, 113, 92, 82, 78, 75, 73, 73, 72,
	  69]
y2 = [118, 106, 78, 72, 62, 55, 47, 43, 39, 37, 35, 30, 27, 26, 21, 22, 20, 21, 18, 17, 17, 15, 14, 14, 12, 12, 12]
y3 = np.repeat(7, 27)

# Set global title.
fig = plt.figure()
st = fig.suptitle("", fontsize="x-large")

# Set tick direction.
matplotlib.rcParams['xtick.direction'] = 'in'
matplotlib.rcParams['ytick.direction'] = 'in'

plt.plot(x1, y1, 'r.-', label='BRaft')
plt.xlabel('Number of nodes (.)')
plt.ylabel('Throughput (tps)')
plt.legend()

plt.plot(x2, y2, 'g.--', label='PBFT')
plt.legend()

plt.plot(x2, y3, 'b.:', label='PoW')
plt.legend()

plt.show()
