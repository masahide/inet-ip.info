# inet-ip.info

## make

```
curl http://geolite.maxmind.com/download/geoip/database/GeoLiteCountry/GeoIP.dat.gz |gunzip >data/GeoIP.dat
go get github.com/jteeuwen/go-bindata
go generate
go get github.com/mitchellh/gox
gox -osarch="linux/386"
```
