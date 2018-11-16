APP?=klyn-test

clean:
rm -f ${APP}

build: clean
go build -o ${APP}

run: build
./${APP}