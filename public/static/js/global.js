"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
// Type assertions for all elements
const thModal = document.getElementById("th-modal");
const thCloseBtn = document.getElementById("close-th-btn");
const thAdd = document.getElementById("add-th-btn");
const replyBtn = document.getElementById("create-post-btn");
const threadCreateBtn = document.getElementById("create-thread-btn");
const thParent = document.getElementById("th-parent");
const signupBtn = document.getElementById("signup-btn");
const responseBoxError = document.getElementById("response-box-error");
const responseMsgError = document.getElementById("resnponse-msg-error");
const responseBoxSuccess = document.getElementById("response-box-success");
const responseMsgSuccess = document.getElementById("resnponse-msg-success");
const signInBtn = document.getElementById("sginin-btn");
if (thModal && thCloseBtn && thAdd) {
    thAdd.addEventListener("click", () => {
        thModal.style.display = "flex";
    });
    thCloseBtn.addEventListener("click", (e) => {
        e.preventDefault();
        thModal.style.display = "none";
    });
}
if (threadCreateBtn) {
    threadCreateBtn.onclick = () => __awaiter(void 0, void 0, void 0, function* () {
        const threadTitleInput = document.getElementById("thread-title");
        const threadContentInput = document.getElementById("thread-content");
        if (!threadTitleInput || !threadContentInput)
            return;
        const threadTitle = threadTitleInput.value;
        const threadContent = threadContentInput.value;
        const response = yield fetch("/create_thread", {
            method: "POST",
            headers: {
                "Content-type": "application/json",
            },
            body: JSON.stringify({
                thread_title: threadTitle,
                thread_content: threadContent,
            }),
        });
        const datas = yield response.json();
        console.log(datas);
        if (response.status === 201) {
            window.location.reload();
        }
    });
}
if (replyBtn && thParent) {
    replyBtn.onclick = () => __awaiter(void 0, void 0, void 0, function* () {
        const postTitleInput = document.getElementById("post-title");
        const postContentInput = document.getElementById("post-content");
        if (!postTitleInput || !postContentInput)
            return;
        const postTitle = postTitleInput.value;
        const postContent = postContentInput.value;
        const threadToken = thParent.dataset.threadTk;
        const response = yield fetch("/create_post", {
            method: "POST",
            headers: {
                "Content-type": "application/json",
            },
            body: JSON.stringify({
                post_title: postTitle,
                post_content: postContent,
                thread_token: threadToken,
            }),
        });
        const datas = yield response.json();
        if (response.status === 201) {
            window.location.reload();
        }
    });
}
if (signupBtn && responseBoxError && responseMsgError && responseBoxSuccess && responseMsgSuccess) {
    signupBtn.onclick = () => __awaiter(void 0, void 0, void 0, function* () {
        try {
            const usernameInput = document.getElementById("reg-username");
            const emailInput = document.getElementById("reg-email");
            const passwordInput = document.getElementById("reg-password");
            if (!usernameInput || !emailInput || !passwordInput)
                return;
            const username = usernameInput.value;
            const email_address = emailInput.value;
            const password = passwordInput.value;
            if (!username) {
                responseBoxError.style.display = "flex";
                responseMsgError.innerText = "username is required !!";
                return;
            }
            if (!email_address) {
                responseBoxError.style.display = "flex";
                responseMsgError.innerText = "email address is required !!";
                return;
            }
            if (!password) {
                responseBoxError.style.display = "flex";
                responseMsgError.innerText = "password is required !!";
                return;
            }
            const response = yield fetch("/auth/signup", {
                method: "POST",
                headers: {
                    "Content-type": "application/json",
                },
                body: JSON.stringify({ username, email_address, password }),
            });
            const data = yield response.json();
            if (response.status !== 201) {
                responseBoxSuccess.style.display = "none";
                responseBoxError.style.display = "flex";
                responseMsgError.innerText = data === null || data === void 0 ? void 0 : data.message;
            }
            else {
                responseBoxError.style.display = "none";
                responseBoxSuccess.style.display = "flex";
                responseMsgSuccess.innerText = "You have successfully created an account. You can login now.";
            }
        }
        catch (err) {
            console.error(err);
        }
    });
}
function CloseErrAlert() {
    if (responseBoxError)
        responseBoxError.style.display = "none";
}
function CloseSuccAlert() {
    if (responseBoxSuccess)
        responseBoxSuccess.style.display = "none";
}
if (signInBtn) {
    signInBtn.onclick = () => __awaiter(void 0, void 0, void 0, function* () {
        try {
            const emailInput = document.getElementById("sginin-email");
            const passwordInput = document.getElementById("sginin-password");
            if (!emailInput || !passwordInput)
                return;
            const email_address = emailInput.value;
            const password = passwordInput.value;
            const response = yield fetch("/auth/signin", {
                method: "POST",
                headers: {
                    "Content-type": "application/json",
                },
                body: JSON.stringify({ email_address, password }),
            });
            const data = yield response.json();
            console.log(data);
            if (response.status === 200) {
                window.location.href = "/";
            }
        }
        catch (err) {
            console.error(err);
        }
    });
}
