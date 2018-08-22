/**
 * 显示input 错误信息 文本框显示为红色
 * @param element
 * @param msg
 */
function showInputError(element, msg, msgdiv = "div") {
    $(element).removeClass("is-valid").addClass("is-invalid");
    if( $(element).parent().children(msgdiv).length == 0) {
        var div = document.createElement(msgdiv);
        $(element).parent().append(div);
    }
    $(element).parent().children(msgdiv).removeClass("valid-feedback").addClass("invalid-feedback").html(msg);
}

/**
 * 显示 input 成功信息  文本框为绿色
 * @param element
 * @param msg
 */
function showInputSuccess(element, msg, msgdiv = "div") {
    $(element).removeClass("is-invalid").addClass("is-valid");
    if( $(element).parent().children(msgdiv).length == 0) {
        var div = document.createElement(msgdiv);
        $(element).parent().append(div);
    }
    $(element).parent().children(msgdiv).removeClass("invalid-feedback").addClass("valid-feedback").html(msg);
}

/**
 * 取消 input 信息框 默认颜色
 * @param element
 * @param msg
 */
function showInputDefault(element, msg, msgdiv = "div") {
    $(element).removeClass("is-invalid").removeClass("is-valid");
    if( $(element).parent().children(msgdiv).length > 0) {
        $(element).parent().children(msgdiv).removeClass("invalid-feedback").removeClass("valid-feedback").html(msg);
    }
}

/**
 * 密码强度检测, 分为0,1,2,3,4 0表示含有非法字符，1弱， 2中， 3强, 4完美
 * @param password
 */
function passwordStrong(password) {
    var level = 0;
    if( /[A-Z]+/.test(password) ) {
        //console.log("password match A-Z")
        level++;
    }
    if( /[a-z]+/.test(password) ) {
        //console.log("password match a-z")
        level++;
    }
    if( /[0-9]+/.test(password) ) {
        //console.log("password match 0-9")
        level++;
    }
    if( /[-_@]+/.test(password) ) {
        //console.log("password match -_@")
        level++;
    }
    if( /[\~\`\$\%\^\&\*\(\)\;\'\"\,\.\?\<\>\+\=]+/.test(password) ) {
        //console.log("password match other")
        level = 0;
    }
    return level;
}