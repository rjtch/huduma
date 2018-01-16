#start from  an image with the latest version of go installed
#and workspace (GOPATH) configured at /project_go, the built context.
FROM golang:1.9.2-alpine3.6 As build

#tools required to build the project.
#`docker build --no-cache` is used to update those dependencies.
RUN apk add --no-cache git
RUN go get github.com/rjtch/huduma

#Gopkg.toml and Gopkg.lock list all project dependencies.
#This layer will only be re-built and reused when Gopkg files are updated
COPY Gopkg.lock Gopkg.toml /project_go/src/github.com/huduma
WORKDIR /project_go/src/github.com/huduma

#install library dependencies
RUN dep ensure -vendor-only

#copy all projects and build it. This layer will be rebuilt 
#when ever a file has changed in the project directory
COPY . /project_go/src/github.com/huduma
RUN go build -o /bin/huduma

FROM alpine
COPY --from=build /bin/huduma /bin/
EXPOSE 3000 3000
ENTRYPOINT [ "/bin/huduma" ]
