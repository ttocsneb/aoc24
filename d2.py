def is_decinc(diffs: list[int]):
    pos = False
    neg = False
    for v in diffs:
        if v < 0:
            if pos:
                return False
            if not pos and not neg:
                neg = True
        else:
            if neg:
                return False
            if not pos and not neg:
                pos = True
    return True


def is_safe(vals: list[int]):
    if len(vals) == 0:
        return False
    offset = vals[:-1]
    diff = [b - a for a, b in zip(vals[1:], offset)]
    abs_diff = [abs(a) for a in diff]

    if len(diff) == 0:
        return False

    if max(abs_diff) > 3 or min(abs_diff) == 0 or not is_decinc(diff):
        return False
    else:
        return True


def gen_alts(vals: list[int]):
    for i in range(len(vals)):
        yield vals[:i] + vals[i + 1 :]


with open("input-d2") as f:
    count = 0
    p2Count = 0
    for line in f.readlines():
        vals = [int(v) for v in line.split() if v]
        if len(vals) == 0:
            continue

        if is_safe(vals):
            count += 1
        else:
            for alt in gen_alts(vals):
                if is_safe(alt):
                    p2Count += 1
                    break

    print(count)
    print(count + p2Count)
