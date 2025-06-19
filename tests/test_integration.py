"""Integration tests for prgen."""

from prgen.llm_client import generate_pr_description


class FakeLLMClient:
    """A simple fake LLM client for testing."""

    def __init__(self, response_title="Test PR Title", response_body="Test PR Body"):
        self.response_title = response_title
        self.response_body = response_body

    def generate_pr_description(self, prompt):
        from prgen.llm_providers import PRResponse

        return PRResponse(title=self.response_title, body=self.response_body)

    def is_available(self):
        return True


class TestPrgenIntegration:
    """Integration tests for prgen workflow."""

    def test_generate_pr_description_happy_path(self, monkeypatch):
        """Test PR description generation with fake LLM client."""

        # Replace the factory with our fake
        def fake_create_client(provider=None):
            return FakeLLMClient("Add new feature", "This PR adds an awesome new feature.")

        monkeypatch.setattr("prgen.llm_client.LLMClientFactory.create_client", fake_create_client)

        # Test with sample diff
        sample_diff = """diff --git a/test.py b/test.py
index 1234567..abcdefg 100644
--- a/test.py
+++ b/test.py
@@ -1,3 +1,6 @@
 def old_function():
     return 'old'
+
+def new_function():
+    return 'hello'
"""

        result = generate_pr_description(sample_diff)

        # Verify the result format
        assert "Add new feature" in result
        assert "This PR adds an awesome new feature." in result

    def test_generate_pr_description_different_response(self, monkeypatch):
        """Test with different PR content."""

        def fake_create_client(provider=None):
            return FakeLLMClient("Fix bug in parser", "Fixed null pointer exception in the parser module.")

        monkeypatch.setattr("prgen.llm_client.LLMClientFactory.create_client", fake_create_client)

        sample_diff = "diff --git a/parser.py b/parser.py\n-    return None\n+    return result"

        result = generate_pr_description(sample_diff)

        assert "Fix bug in parser" in result
        assert "Fixed null pointer exception" in result

    def test_generate_pr_description_error_handling(self, monkeypatch):
        """Test error handling when LLM client creation fails."""

        def fake_create_client(provider=None):
            raise ValueError("No API key available")

        monkeypatch.setattr("prgen.llm_client.LLMClientFactory.create_client", fake_create_client)

        sample_diff = "diff --git a/test.py b/test.py\n+def new_function():\n+    return 'hello'"

        # Should raise the error
        try:
            generate_pr_description(sample_diff)
            raise AssertionError("Expected ValueError to be raised")
        except ValueError as e:
            assert "No API key available" in str(e)

    def test_generate_pr_description_prompt_generation(self, monkeypatch):
        """Test that diff content is properly passed to the LLM."""
        received_prompts = []

        class PromptCapturingClient:
            def generate_pr_description(self, prompt):
                received_prompts.append(prompt)
                from prgen.llm_providers import PRResponse

                return PRResponse(title="Captured", body="Prompt captured successfully")

            def is_available(self):
                return True

        def fake_create_client(provider=None):
            return PromptCapturingClient()

        monkeypatch.setattr("prgen.llm_client.LLMClientFactory.create_client", fake_create_client)

        sample_diff = "diff --git a/important.py b/important.py\n+    important_change = True"

        generate_pr_description(sample_diff)

        # Verify the prompt contains our diff
        assert len(received_prompts) == 1
        assert "important_change = True" in received_prompts[0]
