import requests
import uuid

BASE = "http://localhost:8080/api/auth"

def new_user():
    unique = uuid.uuid4().hex[:8]
    return {
        "username": f"user_{unique}",
        "email": f"user_{unique}@example.com",
        "password": "testpassword123",
    }

def test_full_flow():
    user = new_user()
    s = requests.Session()

    r = s.post(f"{BASE}/signup", json=user)
    assert r.status_code == 201, f"signup failed: {r.status_code} {r.text}"
    print(f"✓ signup: {user['username']}")

    r = s.post(f"{BASE}/login", json={"email": user["email"], "password": user["password"]})
    assert r.status_code == 200, f"login failed: {r.status_code} {r.text}"
    print("✓ login")

    r = s.get(f"{BASE}/test")
    assert r.status_code == 200, f"authenticated request failed: {r.status_code} {r.text}"
    assert user["username"] in r.text, f"expected username in response, got: {r.text}"
    print(f"✓ authenticated: {r.text.strip()}")

    r = s.post(f"{BASE}/logout")
    assert r.status_code == 200, f"logout failed: {r.status_code} {r.text}"
    print("✓ logout")

    r = s.get(f"{BASE}/test")
    assert r.status_code == 401, f"expected 401 after logout, got: {r.status_code}"
    print("✓ rejected after logout")

def test_wrong_password():
    user = new_user()
    s = requests.Session()

    s.post(f"{BASE}/signup", json=user)
    r = s.post(f"{BASE}/login", json={"email": user["email"], "password": "wrongpassword"})
    assert r.status_code == 401, f"expected 401, got: {r.status_code} {r.text}"
    print("✓ rejected wrong password")

def test_no_cookie():
    r = requests.get(f"{BASE}/test")
    assert r.status_code == 401, f"expected 401, got: {r.status_code}"
    print("✓ rejected unauthenticated request")

def test_duplicate_signup():
    user = new_user()
    s = requests.Session()

    s.post(f"{BASE}/signup", json=user)
    r = s.post(f"{BASE}/signup", json=user)
    assert r.status_code == 400, f"expected 400 for duplicate, got: {r.status_code} {r.text}"
    print("✓ rejected duplicate signup")

if __name__ == "__main__":
    tests = [test_no_cookie, test_full_flow, test_wrong_password, test_duplicate_signup]
    passed, failed = 0, 0
    for t in tests:
        try:
            t()
            passed += 1
        except AssertionError as e:
            print(f"✗ {t.__name__}: {e}")
            failed += 1
    print(f"\n{passed} passed, {failed} failed")
