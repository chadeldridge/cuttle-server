#!/usr/bin/perl

use strict;
use JSON;

my $file = shift @ARGV;

$file = "package.json" unless $file;
#die("No package file specified!\nUsage: $0 </path/to/package.json>\n\n") unless $file;
die("Package file not found!\nUsage: $0 </path/to/package.json>\n$!\n") unless -f $file;

my $data = `cat $file`;
$data = decode_json($data);

my $packages;
foreach my $p (sort keys %{$data->{dependencies}}) {
    $packages .= "$p ";
}
system("npm install $packages --save");

$packages = '';
foreach my $p (sort keys %{$data->{devDependencies}}) {
    $packages .= "$p ";
}
system("npm install $packages --save-dev");
