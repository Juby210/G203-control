import QtQuick 2.12
import QtQuick.Layouts 1.12
import QtQuick.Controls 2.14
import QtQuick.Controls 1.6 as QQC1
import QtQuick.Dialogs 1.3
import Backend 1.0

Rectangle {
    SystemPalette { id: palette; colorGroup: SystemPalette.Active }

    id: rectangle
    width: 400
    height: 300
    color: palette.window

    Backend {
        id: backend

        onChangeSearch: s => {
            mainLayout.visible = !s
            searchLayout.visible = s
        }
        onChangeDPI: (d1, d2, d3, d4, d5) => {
            dpi1.text = d1.toString()
            dpi2.text = d2.toString()
            dpi3.text = d3.toString()
            dpi4.text = d4.toString()
            dpi5.text = d5.toString()
        }
    }

    ColumnLayout {
        id: mainLayout
        spacing: 0
        anchors.bottomMargin: 8
        anchors.topMargin: 8
        anchors.leftMargin: 8
        anchors.rightMargin: 8
        anchors.fill: parent
        visible: false

        Button {
            id: selectButton
            Layout.alignment: Qt.AlignHCenter | Qt.AlignTop
            Layout.fillWidth: true

            contentItem: Text {
                text: "Select color"
                color: palette.buttonText
                horizontalAlignment: Text.AlignHCenter
                verticalAlignment: Text.AlignVCenter
            }
            onClicked: {
                colorDialog.visible = true
            }
        }

        GroupBox {
            Layout.alignment: Qt.AlignLeft | Qt.AlignTop
            Layout.fillWidth: true
            title: "Speed"

            Slider {
                id: speedSlider
                anchors.fill: parent
                Layout.fillWidth: true
                Layout.alignment: Qt.AlignHCenter | Qt.AlignTop

                enabled: false
                from: 500
                value: 5000
                to: 12500
            }
        }

        GroupBox {
            id: groupBox
            Layout.alignment: Qt.AlignLeft | Qt.AlignTop
            Layout.fillWidth: true
            title: "Effect"

            RowLayout {
                anchors.fill: parent
                Layout.fillWidth: true
                Layout.alignment: Qt.AlignHCenter | Qt.AlignTop

                RadioButton {
                    id: staticRadio
                    checked: true
                    Layout.fillWidth: true

                    contentItem: Text {
                        text: "Static"
                        color: palette.buttonText
                        leftPadding: staticRadio.indicator.width + staticRadio.spacing
                        verticalAlignment: Text.AlignVCenter
                    }
                    onClicked: {
                        selectButton.enabled = true
                        speedSlider.enabled = false
                    }
                }

                RadioButton {
                    id: breatheRadio
                    Layout.fillWidth: true

                    contentItem: Text {
                        text: "Breathe"
                        color: palette.buttonText
                        leftPadding: staticRadio.indicator.width + staticRadio.spacing
                        verticalAlignment: Text.AlignVCenter
                    }
                    onClicked: {
                        selectButton.enabled = true
                        speedSlider.enabled = true
                    }
                }

                RadioButton {
                    id: cycleRadio
                    Layout.fillWidth: true

                    contentItem: Text {
                        text: "Cycle"
                        color: palette.buttonText
                        leftPadding: staticRadio.indicator.width + staticRadio.spacing
                        verticalAlignment: Text.AlignVCenter
                    }
                    onClicked: {
                        selectButton.enabled = false
                        speedSlider.enabled = true
                    }
                }
            }
        }

        Button {
            Layout.alignment: Qt.AlignHCenter | Qt.AlignTop
            Layout.fillWidth: true

            contentItem: Text {
                text: "Set color"
                color: palette.buttonText
                horizontalAlignment: Text.AlignHCenter
                verticalAlignment: Text.AlignVCenter
            }
            onClicked: {
                let effect = 0
                if (breatheRadio.checked) effect = 1
                else if (cycleRadio.checked) effect = 2
                backend.setColor(selectButton.contentItem.color, speedSlider.value, effect)
            }
        }

        RowLayout {
            id: rowLayout
            Layout.minimumHeight: 50
            Layout.alignment: Qt.AlignLeft | Qt.AlignTop
            Layout.fillWidth: true
            Layout.fillHeight: true

            QQC1.TextField {
                id: dpi1
                textColor: palette.text
                Layout.alignment: Qt.AlignLeft | Qt.AlignBottom
                Layout.fillWidth: true
                validator: IntValidator {
                    bottom: 1
                    top: 8000
                }
            }

            QQC1.TextField {
                id: dpi2
                textColor: palette.text
                Layout.alignment: Qt.AlignLeft | Qt.AlignBottom
                Layout.fillWidth: true
                validator: IntValidator {
                    bottom: 1
                    top: 8000
                }
            }

            QQC1.TextField {
                id: dpi3
                textColor: palette.text
                Layout.alignment: Qt.AlignLeft | Qt.AlignBottom
                Layout.fillWidth: true
                validator: IntValidator {
                    bottom: 1
                    top: 8000
                }
            }

            QQC1.TextField {
                id: dpi4
                textColor: palette.text
                Layout.alignment: Qt.AlignLeft | Qt.AlignBottom
                Layout.fillWidth: true
                validator: IntValidator {
                    bottom: 1
                    top: 8000
                }
            }

            QQC1.TextField {
                id: dpi5
                textColor: palette.text
                Layout.alignment: Qt.AlignLeft | Qt.AlignBottom
                Layout.fillWidth: true
                validator: IntValidator {
                    bottom: 0
                    top: 8000
                }
            }
        }

        Button {
            Layout.topMargin: 1
            Layout.alignment: Qt.AlignHCenter | Qt.AlignTop
            Layout.fillWidth: true

            contentItem: Text {
                text: "Set DPI"
                color: palette.buttonText
                horizontalAlignment: Text.AlignHCenter
                verticalAlignment: Text.AlignVCenter
            }
            onClicked: {
                backend.setDPI(Number(dpi1.text), Number(dpi2.text), Number(dpi3.text), Number(dpi4.text), Number(dpi5.text))
            }
        }
    }

    ColumnLayout {
        id: searchLayout
        anchors.fill: parent
        visible: true

        BusyIndicator {
            id: busyIndicator
            Layout.preferredHeight: 120
            Layout.preferredWidth: 120
            Layout.alignment: Qt.AlignHCenter | Qt.AlignVCenter
        }

        Text {
            id: searchText

            text: "Searching for G203 mouse"
            Layout.alignment: Qt.AlignHCenter | Qt.AlignTop
            font.pixelSize: 20
        }
    }

    ColorDialog {
        id: colorDialog

        title: "Please choose a color"
        onAccepted: {
            console.log("ColorDialog: " + colorDialog.color)
            selectButton.contentItem.color = colorDialog.color
        }
    }
}
