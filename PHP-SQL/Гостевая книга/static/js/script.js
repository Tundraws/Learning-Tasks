document.getElementById('message-form').addEventListener('submit', function(e) {
    const content = document.getElementById('content').value.trim();
    
    if (content === '') {
        e.preventDefault();
        alert('Сообщение не может быть пустым');
    }
});

function addMessageToTop(message) {
    const messagesDiv = document.querySelector('.messages');
    const messageHtml = `
        <div class="message">
            <div class="message-header">
                <span class="date">${message.date}</span>
                <span class="name">${message.name}</span>
            </div>
            <div class="message-content">${message.content}</div>
        </div>
    `;
    messagesDiv.insertAdjacentHTML('afterbegin', messageHtml);
}