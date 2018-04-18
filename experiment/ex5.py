import matplotlib
import matplotlib.pyplot as plt
import numpy as np

x1 = np.linspace(6, 30, 25)
x2 = np.linspace(4, 30, 27)

# BRaft
y1 = [18.8, 15.2, 16.5, 17.0, 12.6, 15.0, 13.8, 15.1, 10.21, 13.8, 13.7, 12.8, 15.1, 15.3, 10.3, 14.4, 15.2, 15.5, 13.9,
	  10.6, 13.3, 14.8, 11.5, 14.1, 10.0]

y2 = [16.3, 15.7, 11.5, 12.3, 14.8, 13.1, 13.6, 12.6, 13.2, 12.8, 14.1, 17.4, 12.6, 15.9, 13.6, 12.9, 12.3, 17.7, 11.9,
	  11.8, 15.6, 19.7, 12.8, 11.5, 12.3, 13.2, 15.2]

# Set global title.
fig = plt.figure()
st = fig.suptitle("", fontsize="x-large")

# Set tick direction.
matplotlib.rcParams['xtick.direction'] = 'in'
matplotlib.rcParams['ytick.direction'] = 'in'

plt.plot(x1, y1, 'r.-', label='BRaft')
plt.xlabel('Number of nodes (.)')
plt.ylabel('Latency (ms)')
plt.legend()

plt.plot(x2, y2, 'g.--', label='PBFT')
plt.legend()

plt.show()
