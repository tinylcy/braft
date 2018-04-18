import matplotlib
import matplotlib.pyplot as plt
import matplotlib.patches as patches
import matplotlib.path as path
import pylab

import numpy as np

x = np.linspace(6, 30, 25)

y = [1, 2, 2, 1, 2, 1, 3, 3, 3, 3, 4, 4, 3, 4, 5, 6, 4, 4, 6, 4, 5, 5, 4, 6, 6]

# Set tick direction.
matplotlib.rcParams['xtick.direction'] = 'in'
matplotlib.rcParams['ytick.direction'] = 'in'

plt.plot(x, y, '.-')
plt.xlabel('Number of nodes (.)')
plt.ylabel('Time used (ms)')

plt.show()
