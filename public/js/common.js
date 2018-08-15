/**
 * 显示input 错误信息 文本框显示为红色
 * @param element
 * @param msg
 */
function show_input_error(element, msg, msgdiv = "div") {
    $(element).removeClass("is-valid").addClass("is-invalid");
    if( $(element).parent().children(msgdiv).length == 0) {
        var div = document.createElement(msgdiv)
        $(element).parent().append(div)
    }
    $(element).parent().children(msgdiv).removeClass("valid-feedback").addClass("invalid-feedback").html(msg)
}

/**
 * 显示 input 成功信息  文本框为绿色
 * @param element
 * @param msg
 */
function show_input_success(element, msg, msgdiv = "div") {
    $(element).removeClass("is-invalid").addClass("is-valid");
    if( $(element).parent().children(msgdiv).length == 0) {
        var div = document.createElement(msgdiv)
        $(element).parent().append(div)
    }
    $(element).parent().children(msgdiv).removeClass("invalid-feedback").addClass("valid-feedback").html(msg)
}

/**
 * 取消 input 信息框 默认颜色
 * @param element
 * @param msg
 */
function show_input_default(element, msg, msgdiv = "div") {
    $(element).removeClass("is-invalid").removeClass("is-valid")
    if( $(element).parent().children(msgdiv).length > 0) {
        $(element).parent().children(msgdiv).removeClass("invalid-feedback").removeClass("valid-feedback").html(msg)
    }

}