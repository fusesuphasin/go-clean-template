FROM mongo:3.1
#FROM : อ้างอิงถึง parent image ชื่อ mongo ติด tag version 3.4.7 

ENV MONGO_INITDB_ROOT_USERNAME admin
ENV MONGO_INITDB_ROOT_PASSWORD 1234
ENV MONGO_INITDB_DATABASE myLearning
# กำหนด ENV ให้มัน 

#กรณีที่ต้องการ Set initial db ให้มัน
ADD ./data/datastudent.js /docker-entrypoint-initdb.d


#ดู docs แบบ official ได้ที่ https://hub.docker.com/_/mongo