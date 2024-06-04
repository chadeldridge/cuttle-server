FROM ubuntu:latest

RUN apt update && apt install openssh-server sudo -y

RUN useradd -rm -d /home/bob -s /bin/bash -g root -G sudo -u 1001 bob
RUN echo 'bob:testUserP@ssw0rd' | chpasswd

RUN mkdir /home/bob/.ssh
RUN chmod 700 /home/bob/.ssh
RUN chown bob:root /home/bob/.ssh/

COPY keys /home/bob/.ssh/authorized_keys
RUN chown bob:root /home/bob/.ssh/authorized_keys

RUN service ssh start

EXPOSE 22

CMD ["/usr/sbin/sshd","-D"]
