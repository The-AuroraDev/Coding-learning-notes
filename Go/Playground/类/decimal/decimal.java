public class Coordinate {
    private final double d, m, s;
    private final char h;

    public Coordinate(double d, double m, double s, char h) {
        this.d = d;
        this.m = m;
        this.s = s;
        this.h = h;
    }

    // DMS -> DD
    public double decimal() {
        double sign = 1.0;
        switch (h) {
            case 'S': case 'W': case 's': case 'w':
                sign = -1;
                break;
        }
        return sign * (d + m / 60 + s / 3600);
    }

    public static void main(String[] args) {
        Coordinate lat = new Coordinate(30, 39, 34.2, 'N');
        Coordinate lng = new Coordinate(104, 3, 56.52, 'E');  // 不能叫 long！

        System.out.println(lat.decimal() + " " + lng.decimal());
    }
}