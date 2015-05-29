This is a Go library and a collection of tools to interact with and analyze Tor HSDirs.

## Bruteforcing tool

`brute` will bruteforce Identity Keys that will be the 6 HSDir for the given onion
at the given time, considering the given consensus state (which should probably just be the most
[recent](https://collector.torproject.org/recent/relay-descriptors/consensuses/).

    ./bin/brute 2015-05-23-23-00-00-consensus facebookcorewwwi.onion 2015-05-29T12:00:00Z

You should really patch your Go standard library with the provided reduced-primalty patch, it will make the bruteforce
20 times as fast and complete checks are performed anyway by the tool itself (ProTip: have a separate Go tree for this
security-shattering patch!)

## Preprocessing

First, download consensus files from the following two addresses:

https://collector.torproject.org/archive/relay-descriptors/consensuses/  
https://collector.torproject.org/recent/relay-descriptors/consensuses/

and make sure they are named like this (this is what you will get by just unpacking the tarballs):

    $DIR/consensuses-2015-01/01/2015-01-01-00-00-00-consensus

Then there are two steps to preprocessing:

    preprocess $DIR 2015-01-01-00 2015-05-28-10
    grind pckcns.dat > stats.jsonl

This should leave you with the following three files: `pckcns.dat`, `keys.db` and `stats.jsonl`

## Analysis tools

### scrolls

`scrolls` will go through consensuses in a time frame and print the HSDir set every time it changed,
with color-coded suspiciousness scores.

    ./bin/scrolls pckcns.dat keys.db stats.jsonl facebookcorewwwi.onion 2015-05-01-01 2015-05-28-10
    
![image](https://cloud.githubusercontent.com/assets/1225294/7873542/e206bdf6-05a3-11e5-99db-d9d89a2bf4ec.png)

* "Dist score" is a normalized score of how near the first HSDir is to the Descriptor ID, compared to the
  expected mean and average deviation on the ring.
* "Dist4 score" is the same, but for the distance of the 4th node.
* "Age" is how many hours before the node got the HSDir flag (or appeared). ∞ means since the beginning of the time frame.
* "Long" is how many hours after the node lost the HSDir flag (or disappeared). ∞ means until the end of the time frame.
* "Colo keys" is how many other Identity Keys have been observed on IPs that hosted this Identity Key.

### curiosity

`curiosity` will generate statistics about correlations between IPs and multiple Identity Keys.

    ./bin/curiosity data/keys.db colocated | sort -g
    
    [...]
    57 keys on 1 IPs - EFA4BF05D097357832B8C7EE3391C8A534CB06F1 [178.32.143.175]
    57 keys on 1 IPs - FADD1E620505CD8B1E22DC9CC503D06592C24EF9 [178.32.143.175]
    57 keys on 1 IPs - FFCCB39F28DB9926F9D34D97F1D98EEA834EED38 [178.32.143.175]
    57 keys on 3 IPs - 0045D088EB8E8E3634F3C82A8010A50BC18938D4 [54.242.144.249 23.20.33.81 23.22.192.88]
    57 keys on 3 IPs - 01120D3FA96BA5B5F5FBB5E1581E8C73E0EFE24D [54.242.144.249 23.20.33.81 23.22.192.88]
    57 keys on 3 IPs - 01C5B224668F77D3662956DDEC15DA46A252CC23 [54.242.144.249 23.20.33.81 23.22.192.88]
    57 keys on 3 IPs - 027B74EEB1A43CFA490D3AE7C4CCC11F5F26A3A7 [54.242.144.249 23.20.33.81 23.22.192.88]
    57 keys on 3 IPs - 03D401D7B671098945BD5B14E3D54278CBCF5067 [54.242.144.249 23.20.33.81 23.22.192.88]

## Coming soon

* automated orchestration tool to easily deploy defensive arrays of HSDirs
* real-time targeted and generic statistical monitoring tools
