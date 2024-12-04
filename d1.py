with open("input-d1") as f:
    total = 0

    a = []
    b = []
    for line in f.readlines():
        foo = [int(a.strip()) for a in line.split()]

        if not foo:
            continue

        a.append(foo[0])
        b.append(foo[1])

    a.sort()
    b.sort()

    print(sum([abs(j - i) for i, j in zip(a, b)]))

    a_c = {}
    b_c = {}

    for v in a:
        a_c.setdefault(v, 0)
        a_c[v] += 1
    for v in b:
        b_c.setdefault(v, 0)
        b_c[v] += 1

    sim = 0
    for k, v in a_c.items():
        sim += k * b_c.get(k, 0) * v

    print(sim)
