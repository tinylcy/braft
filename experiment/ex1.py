import matplotlib
import matplotlib.pyplot as plt
import numpy as np

x = np.linspace(6, 30, 25)

y1 = [501, 501, 502, 501, 501, 502, 502, 502, 502, 502, 502, 502, 503, 502, 503, 503, 503, 504, 503, 504, 504, 504, 505,
	  505, 505]
y2 = [502, 502, 501, 502, 501, 501, 502, 502, 502, 502, 504, 503, 503, 504, 503, 503, 504, 504, 504, 504, 504, 504, 505,
	  505, 505]

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
