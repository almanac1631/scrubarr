import {Ref, ref} from "vue";

export const NotificationType = {
    Info: 'info',
    Success: 'success',
    Error: 'error',
} as const;
type NotificationTypeEnum = typeof NotificationType[keyof typeof NotificationType];

export interface NotificationEntry {
    message: string;
    type: NotificationTypeEnum;
}

export const notificationList: Ref<Array<NotificationEntry>> = ref([]);

export function notify(message: string, type: NotificationTypeEnum): void {
    const notification = {
        message: message,
        type: type,
    }
    notificationList.value.push(notification);
    setTimeout(() => {
        notificationList.value.splice(notificationList.value.indexOf(notification), 1);
    }, 3000);
}
