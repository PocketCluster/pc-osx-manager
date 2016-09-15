#!/usr/bin/env python

a = [i for i in xrange(0, 100)]
print a
print "-" * 80
#cg = []; map(lambda e, n: cg.append(n * 3) if (n % 3 == 0) else cg.append(n * 0), a); print cg
#print reduce(lambda l, r: l + r, a)

def redfunc(e, n):
    if not e:
        cg = list()
        cg.append(n)
        return cg
    else:
        e.append(n)
#print reduce(redfunc, a)
#print reduce(lambda l, r: [[r]] if not l else l.append([r]), a)
