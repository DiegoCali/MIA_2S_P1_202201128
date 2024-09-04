export async function POST<T>(path: string, data: T) {
    console.log(data);
    const res = await fetch(path, {
        method: 'POST',
        body: JSON.stringify(data),
        headers: {
            'Content-Type': 'application/json',
        },
    });
    if (!res.ok) {
        throw new Error(res.statusText);
    }
    return res.json();
}