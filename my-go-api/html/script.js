document.addEventListener('DOMContentLoaded', () => {
    const form = document.getElementById('bookForm');
    const bookList = document.getElementById('bookList');
    const resetButton = document.getElementById('resetButton');

    form.addEventListener('submit', async (event) => {
        event.preventDefault();

        const newBook = {
            title: document.getElementById('title').value,
            author: document.getElementById('author').value,
            genre: document.getElementById('genre').value,
            description: document.getElementById('description').value,
            isbn: document.getElementById('isbn').value,
            image: document.getElementById('image').value,
            published: document.getElementById('published').value,
            publisher: document.getElementById('publisher').value
        };

        try {
            const response = await fetch('/add', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(newBook)
            });

            if (response.ok) {
                const addedBook = await response.json();
                addBookToList(addedBook);
                form.reset();
            } else {
                console.error('Failed to add book');
            }
        } catch (error) {
            console.error('Error', error);
        }
    });

    resetButton.addEventListener('click', async () => {
        try{
            const response = await fetch('/reset', {
                method: 'POST'
            });

            if (response.ok){
                bookList.innerHTML = '';
                location.reload();
            }
            else {
                console.error('Failed to reset books');
            }
        }
        catch (error){
            console.error('Error', error);
        }
    });

    function addBookToList(book) {
        const li = document.createElement('li');
        li.classList.add('book-item');
        li.setAttribute('data-id', book.id);
        li.innerHTML = `
            <h3>${book.title}</h3>
            <p>Author: ${book.author}</p>
            <p>Genre: ${book.genre}</p>
            <p>${book.description}</p>
            <p>ISBN: ${book.isbn}</p>
            <img src="${book.image}" alt="${book.title}">
            <p>Published: ${book.published}</p>
            <p>Publisher: ${book.publisher}</p>
            <div class="book-actions">
                <button class="edit-btn">Edit</button>
                <button class="delete-btn">Delete</button>
            </div>
        `;
        bookList.appendChild(li);
    }

    bookList.addEventListener('click', async (event) => {
        const target = event.target;
        const li = target.closest('li.book-item');
        const bookId = li ? li.getAttribute('data-id') : null;

        if (target.classList.contains('edit-btn') && bookId) {
            editBook(parseInt(bookId));
        } else if (target.classList.contains('delete-btn') && bookId) {
            deleteBook(parseInt(bookId));
        }
    });

    async function editBook(id) {
        const title = prompt('Enter new title:');
        if (!title) return;

        const bookToUpdate = {
            title,
            author: prompt('Enter new author:'),
            genre: prompt('Enter new genre:'),
            description: prompt('Enter new description:'),
            isbn: prompt('Enter new ISBN:'),
            image: prompt('Enter new image URL:'),
            published: prompt('Enter new published date:'),
            publisher: prompt('Enter new publisher:')
        };

        try {
            const response = await fetch(`/update?id=${id}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(bookToUpdate)
            });

            if (response.ok) {
                const updatedBook = await response.json();
                updateBookInList(updatedBook);
            } else {
                console.error('Failed to update book');
            }
        } catch (error) {
            console.error('Error', error);
        }
    }

    async function deleteBook(id) {
        if (!confirm('Are you sure you want to delete this book?')) return;

        try {
            const response = await fetch(`/delete?id=${id}`, {
                method: 'DELETE'
            });

            if (response.ok) {
                removeBookFromList(id);
            } else {
                console.error('Failed to delete book');
            }
        } catch (error) {
            console.error('Error', error);
        }
    }

    function updateBookInList(updatedBook) {
        const li = bookList.querySelector(`li[data-id='${updatedBook.id}']`);
        if (li) {
            li.innerHTML = `
                <h3>${updatedBook.title}</h3>
                <p>Author: ${updatedBook.author}</p>
                <p>Genre: ${updatedBook.genre}</p>
                <p>${updatedBook.description}</p>
                <p>ISBN: ${updatedBook.isbn}</p>
                <img src="${updatedBook.image}" alt="${updatedBook.title}">
                <p>Published: ${updatedBook.published}</p>
                <p>Publisher: ${updatedBook.publisher}</p>
                <div class="book-actions">
                    <button class="edit-btn">Edit</button>
                    <button class="delete-btn">Delete</button>
                </div>
            `;
        }
    }

    function removeBookFromList(id) {
        const li = bookList.querySelector(`li[data-id='${id}']`);
        if (li) {
            li.remove();
        }
    }
});
