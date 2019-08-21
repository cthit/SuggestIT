

export const translateTimestamp = (ts) => {
    var pre = "Inlagd "
    var now = new Date();
    var diff = now.getTime() - new Date(ts).getTime() + now.getTimezoneOffset()*60000;
    diff /= 1000;
    if(diff <= 60){
        return pre + parseInt(diff, 10) + " sekunder sedan";
    }
    diff /=60;
    if(diff <= 60){
        return pre + parseInt(diff, 10) + " minuter sedan";
    }
    diff /=60;
    if(diff <= 24){
        return pre + parseInt(diff, 10) + " timmar sedan";
    }
    diff /=24;
    if(diff <= 365.25/12){
        return pre + parseInt(diff, 10) + " dagar sedan";
    }
    diff /= 365.25/12;
    if(diff <= 12){
        return pre + parseInt(diff, 10) + " månader sedan";
    }
    diff /= 12;
    return pre + parseInt(diff, 10) + " år sedan";
}
