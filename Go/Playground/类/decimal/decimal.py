class Coordinate:
    def __init__(self, d, m, s, h):
        self.d = d
        self.m = m
        self.s = s
        self.h = h

    # DMS -> DD
    def decimal(self):
        sign = 1.0
        if self.h in ('S', 'W', 's', 'w'):
            sign = -1.0
        return sign * (self.d + self.m / 60 + self.s / 3600)


if __name__ == '__main__':
    # ChengDu: 30°39′34.2″N 104°03′56.52″E
    lat = Coordinate(30, 39, 34.2, 'N')
    lng = Coordinate(104, 3, 56.52, 'E')

    print(lat.decimal(), lng.decimal())
    # 30.659499999999998 104.06569999999999