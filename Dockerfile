# Copyright 2018 The Knative Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM golang AS builder

WORKDIR /
ADD . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o go4autoscaler ./

FROM ubuntu:18.04
RUN \
      apt-get update && \
      apt-get -qqy install --fix-missing \
            curl \
            ca-certificates \
      && \
      apt-get clean

EXPOSE 8080
COPY --from=builder ./go4autoscaler /go4autoscaler

ENTRYPOINT ["/go4autoscaler"]