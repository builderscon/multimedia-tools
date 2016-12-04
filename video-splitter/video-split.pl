use strict;
use JSON ();

my ($src, $fn) = @ARGV;

my $config = JSON::decode_json(do{ open my $fh, '<', $fn; local $/; <$fh> });

for my $data (@$config) {
    my @cmd = qw(ffmpeg -y);
    my $start_sec = to_sec($data->{start});
    my $end_sec   = to_sec($data->{end});
    my $end       = from_sec($end_sec - $start_sec);
    push @cmd, ('-ss', $data->{start});
    push @cmd, ('-i', $src);
    push @cmd, ('-t', $end);

    if (my $cb_data = $data->{cropblur}) {
        my $filter = sprintf(
            "[0:v]crop=%d:%d:%d:%d,boxblur=luma_radius=15",
            $cb_data->{width}, $cb_data->{height}, $cb_data->{x}, $cb_data->{y}
        );


        if (my $t_data = $cb_data->{between}) {
            $filter .= ":enable='"
            my @between;
            for my $between (@$t_data) {
                push @between, sprintf("between(t,%d,%d)", to_sec($between->{start}), to_sec($between->{end}))
            }
            $filter .= join('+', @between);
            $filter .= "'[fg]";
        }
        $filter .= sprintf("; [0:v][fg]overlay=%d:%d[v]", $cb_data->{x}, $cb_data->{y});

        push @cmd, (
            "-filter_complex", $filter,
            "-map", "[v]",
        )
    } else {
        push @cmd, qw(-c:v copy);
    }
    push @cmd, qw(-filter_complex), "highpass=f=500, lowpass=f=4000, dynaudnorm, volume=2dB";
    push @cmd, $data->{filename};
    system(@cmd);
}

sub to_sec {
    my $s = shift;
    if ($s !~ /^(\d+):(\d+):(\d+)/) {
        die "failed to parse $s";
    }

    return $1 * 3600 + $2 * 60 + $3;
}

sub from_sec {
    my $v = shift;
    my $h = int($v / 3600);
    my $m = int(($v - $h * 3600)/60);
    my $s = $v % 60;
    return sprintf("%02d:%02d:%02d", $h, $m, $s)
}
