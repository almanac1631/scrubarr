import {Ref, ref} from "vue";

export interface NotificationEntry {
    message: string;
}

export const notificationList: Ref<Array<NotificationEntry>> = ref([]);

export function notify(message: string): void {
    const notification = {
        message: message,
    }
    notificationList.value.push(notification);
    setTimeout(() => {
        notificationList.value.splice(notificationList.value.indexOf(notification), 1);
    }, 3000);
}
