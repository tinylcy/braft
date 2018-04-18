import matplotlib
import matplotlib.pyplot as plt
import numpy as np

x = np.linspace(6, 30, 25)

y1 = [61, 65, 55, 56, 59, 66, 65, 70, 57, 51, 76, 56, 50, 61, 51, 58, 51, 52, 50, 49, 53, 55, 56, 52, 49]
y2 = [635, 632, 628, 630, 627, 625, 610, 625, 645, 663, 667, 691, 716, 787, 790, 893, 878, 1082, 1214, 1266, 1316, 1357,
	  1363, 1388, 1429]

# Set global title.
fig = plt.figure()
st = fig.suptitle("", fontsize="x-large")

# Set tick direction.
matplotlib.rcParams['xtick.direction'] = 'in'
matplotlib.rcParams['ytick.direction'] = 'in'

plt.subplot(1, 2, 1)
plt.plot(x, y1, '.-', label='Raft')
plt.xlabel('Number of nodes (.)')
plt.ylabel('Time used (ms)')
plt.legend()

plt.subplot(1, 2, 2)
plt.plot(x, y2, '.-', label='BRaft')
plt.xlabel('Number of nodes (.)')
plt.ylabel('Time used (ms)')
plt.legend()

plt.show()
