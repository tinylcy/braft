import matplotlib
import matplotlib.pyplot as plt
import numpy as np

x = np.linspace(100, 1000, 10)
y1 = [644, 1254, 1867, 2417, 3146, 3802, 4527, 5088, 5769, 6548]
y2 = [66, 151, 247, 389, 516, 651, 841, 1029, 1208, 1392]

# Set global title.
fig = plt.figure()
st = fig.suptitle("", fontsize="x-large")

# Set tick direction.
matplotlib.rcParams['xtick.direction'] = 'in'
matplotlib.rcParams['ytick.direction'] = 'in'

plt.plot(x, y1, 'r.-', label='BRaft')
plt.xlabel('Number of nodes (.)')
plt.ylabel('Time used (ms)')

plt.plot(x, y2, 'g.--', label='Raft')
plt.xlabel('Number of client requests (.)')
plt.ylabel('Time used (ms)')

plt.legend()
plt.show()
