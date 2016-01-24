import QtQuick 2.2
import QtQuick.Controls 1.1
import QtQuick.Layouts 1.0

ApplicationWindow {
    id: applicationWindow
    visible: true
    title: "Gmail Notifications Client"
    property int margin: 50
    width: mainLayout.implicitWidth + 2 * (margin * 3)
    height: mainLayout.implicitHeight + 2 * margin
    minimumWidth: mainLayout.Layout.minimumWidth + 2 * (margin * 3)
    minimumHeight: mainLayout.Layout.minimumHeight + 2 * margin

    ColumnLayout {
        id: mainLayout
        anchors.fill: parent
        anchors.margins: margin

        GroupBox {
            id: rowBox
            title: "Configuration"
            Layout.fillWidth: true

            Column {
                id: rowLayout
                anchors.fill: parent
                property int timeout: timeoutField.text
                property int notificationsLimit: notificationsLimitField.text
                property alias notification: notificationLabel.text
                property alias notificationColor: notificationLabel.color

                Label {
                    text: "Timeout"
                }
                TextField {
                    id: timeoutField
                    placeholderText: "Timeout (s)"
                    text: controller.config.timeout
                    Layout.fillWidth: true
                }

                Label {
                    text: "Notifications Limit"
                }
                TextField {
                    id: notificationsLimitField
                    placeholderText: "Notifications Limit"
                    text: controller.config.notificationsLimit
                    Layout.fillWidth: true
                }
                Button {
                    id: editForm
                    text: "Save"
                    onClicked:
                        controller.saveButtonReleased(parent)
                }

                Text {
                    id: notificationLabel
                    text: controller.notification.text
                    color: controller.notification.color
                    Layout.fillWidth: true
                }

                Button {
                    id: runAgain
                    text: "Reload with new configuration"
                    onClicked:
                        controller.runAgainButtonReleased(parent)
                }
            }
        }
    }
}
