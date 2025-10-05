# pihole-pause

Quickly pause Pi Hole DNS blocking for a duration of time. Default is 60s

## Usage

```
# No duration
pihole-pause

Enter pi-hole password: 
2025/10/05 19:03:08 INFO done response="{\"blocking\":\"disabled\",\"timer\":60,\"took\":0.01505279541015625}"
2025/10/05 19:03:08 INFO done



# Specify duration
pause-pihole 2m
Enter pi-hole password: 
2025/10/05 19:05:15 INFO done response="{\"blocking\":\"disabled\",\"timer\":120,\"took\":5.1736831665039062e-05}"
2025/10/05 19:05:15 INFO done
```

## Install 

```
make install
```

To uninstall it

```
make uninstall
```

