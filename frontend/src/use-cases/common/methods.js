export const translateTimestamp = ts => {
    var pre = "Posted ";
    var now = new Date();
    var diff =
        now.getTime() -
        new Date(ts).getTime() +
        now.getTimezoneOffset() * 60000;
    diff /= 1000;
    if (diff <= 60) {
        return pre + parseInt(diff, 10) + " seconds ago";
    }
    diff /= 60;
    if (diff <= 60) {
        return pre + parseInt(diff, 10) + " minutes ago";
    }
    diff /= 60;
    if (diff <= 24) {
        return pre + parseInt(diff, 10) + " hours ago";
    }
    diff /= 24;
    if (diff <= 365.25 / 12) {
        return pre + parseInt(diff, 10) + " gays ago";
    }
    diff /= 365.25 / 12;
    if (diff <= 12) {
        return pre + parseInt(diff, 10) + " months ago";
    }
    diff /= 12;
    return pre + parseInt(diff, 10) + " years ago";
};
