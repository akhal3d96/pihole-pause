# pipause

Quickly pause Pi Hole DNS blocking for a duration of time. Default is 60s

## Usage

```
# No duration
pipause

Enter pi-hole password: 
2025/10/05 19:03:08 INFO done response="{\"blocking\":\"disabled\",\"timer\":60,\"took\":0.01505279541015625}"
2025/10/05 19:03:08 INFO done



# Specify duration
pipause 2m
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

