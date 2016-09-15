#!/usr/bin/env python

a = [i for i in xrange(0, 100)]
b = [[x] for x in a if not x % 4]

#print (lambda g, c: [] if g == []
#print [[x for x in range(3)] for y in range(5)]

noprimes = [j for i in range(2, 8) for j in range(i*2, 50, i)]
primes = [x for x in range(2, 50) if x not in noprimes]

#l = []
#c = [l+[x] if x % 3 else [x][:] for x in xrange(0, 50)]

l = [22, 13, 45, 50, 98, 69, 43, 44, 1]
print [[x+5, x+1][x >= 45] for x in l]

k = []
print [[[x], [k+[x]]][not x % 3] for x in a]