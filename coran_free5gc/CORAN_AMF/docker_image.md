#### To generate image:
```bash
docker build -t <name_of_image>:<tag> .
```
*Example:* `docker build -t coran_amf:v1.1`

#### To run docker image
```bash
docker images 
docker run --name <name_of_container> <image_name>:<tag>
```
*Example:* `docker run --name coran_amf coran_amf:v1.1`
