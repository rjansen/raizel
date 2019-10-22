FROM google/cloud-sdk:alpine

# Install Java 8 for Datastore emulator
RUN apk add --update --no-cache openjdk8-jre &&\
    gcloud components install cloud-firestore-emulator beta --quiet

# Volume to persist Datastore data
VOLUME /opt/data

EXPOSE 8086

RUN gcloud config set project raizel

ENTRYPOINT ["gcloud", "beta", "emulators", "firestore", "start", "--host-port=0.0.0.0:8086"]
