import matplotlib
import matplotlib.pyplot as plt
import matplotlib.patches as patches
import matplotlib.path as path
import pylab

import numpy as np

x = np.linspace(6, 30, 25)

y1 = [181, 260, 151, 171, 251, 231, 170, 258, 248, 266, 299, 260, 151, 171, 251, 287, 170, 258, 248, 231, 173, 186, 231,
	  258, 266]

y2 = [1, 2, 2, 1, 2, 1, 3, 3, 3, 3, 4, 4, 3, 4, 5, 6, 4, 4, 6, 4, 5, 5, 4, 6, 6]

# Set tick direction.
matplotlib.rcParams['xtick.direction'] = 'in'
matplotlib.rcParams['ytick.direction'] = 'in'

plt.bar(x, y1, label='Phase 1', hatch='/')
plt.bar(x, y2, bottom=y1, label='Phase 2')
plt.xlabel('Number of nodes (.)')
plt.ylabel('Time used (ms)')

plt.legend()
plt.show()
