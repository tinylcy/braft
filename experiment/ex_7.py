import matplotlib
import matplotlib.pyplot as plt
import numpy as np

x = np.linspace(6, 30, 25)

y1 = [1.17, 2.4, 2.09, 2.51, 3.01, 3.18, 4.73, 5.01, 5.82, 5.17, 5.41, 7.69, 7.97, 7.78, 7.52, 9.47, 11.07, 10.92,
	  12.03, 11.97, 13.89, 14.47, 15.16, 14.22, 14.99]

# Set global title.
fig = plt.figure()
st = fig.suptitle("", fontsize="x-large")

# Set tick direction.
matplotlib.rcParams['xtick.direction'] = 'in'
matplotlib.rcParams['ytick.direction'] = 'in'

plt.plot(x, y1, '.-')
plt.xlabel('Number of nodes (.)')
plt.ylabel('Time used (ms)')

plt.show()
