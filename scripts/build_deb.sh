#!/usr/bin/env bash

set -e

# define some important variables
PACKAGE_VERSION=$(cat /workspace/VERSION | cut -c 2-)
PACKAGE_NAME="device-init"
BUILD_DIR="/deb_build"
BINARY_SIZE="$(stat -c %s /workspace/device-init_linux_arm)"
DESCRIPTION="Initialise your device on boot with user defined configuration"

# create build dir where we assemble the final debian package
mkdir -p ${BUILD_DIR}/package/${PACKAGE_NAME}

# copy package template to build directory
cp -r /workspace/package/* ${BUILD_DIR}/package/${PACKAGE_NAME}/

# copy package control template and replace version info
sed -i'' "s/<VERSION>/${PACKAGE_VERSION}/g" ${BUILD_DIR}/package/${PACKAGE_NAME}/DEBIAN/control
sed -i'' "s/<NAME>/${PACKAGE_NAME}/g" ${BUILD_DIR}/package/${PACKAGE_NAME}/DEBIAN/control
sed -i'' "s/<SIZE>/${BINARY_SIZE}/g" ${BUILD_DIR}/package/${PACKAGE_NAME}/DEBIAN/control
sed -i'' "s/<DESCRIPTION>/${DESCRIPTION}/g" ${BUILD_DIR}/package/${PACKAGE_NAME}/DEBIAN/control
sed -i'' "s/<DEPENDS>//g" ${BUILD_DIR}/package/${PACKAGE_NAME}/DEBIAN/control

# copy binary that will be packaged to destination folder
cp /workspace/device-init_linux_arm ${BUILD_DIR}/package/${PACKAGE_NAME}/usr/local/bin
# prevent .gitignore from ending up in the package
rm ${BUILD_DIR}/package/${PACKAGE_NAME}/usr/local/bin/.gitignore

# ensure that the travis-ci user can access the sd-card image file
umask 0000

# create package with dpkg-deb
cd ${BUILD_DIR}/package && dpkg-deb --build ${PACKAGE_NAME}

cp ${BUILD_DIR}/package/${PACKAGE_NAME}.deb /workspace/${PACKAGE_NAME}-${PACKAGE_VERSION}-armhf.deb
chmod 777 /workspace/${PACKAGE_NAME}-${PACKAGE_VERSION}-armhf.deb
