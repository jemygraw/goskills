package goskills

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseAllSkills(t *testing.T) {
	skillsRoot := "examples/skills"
	
	entries, err := os.ReadDir(skillsRoot)
	if err != nil {
		t.Fatalf("Error: Could not read skills directory '%s': %v", skillsRoot, err)
	}

	parsedCount := 0
	var artifactsBuilderSkill *SkillPackage

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		skillPath := filepath.Join(skillsRoot, entry.Name())
		skillPackage, err := ParseSkillPackage(skillPath)
		
		// We expect some directories to fail parsing, so we only log non-fatal errors
		if err != nil {
			// These are not skills, so we expect them to fail.
			if !strings.Contains(entry.Name(), ".git") && !strings.Contains(entry.Name(), ".claude-plugin") && !strings.Contains(entry.Name(), "document-skills") {
				t.Errorf("Parsing skill '%s' failed unexpectedly: %v", entry.Name(), err)
			}
			continue
		}

		if skillPackage.Meta.Name == "artifacts-builder" {
			artifactsBuilderSkill = skillPackage
		}

		parsedCount++
	}

	// Assertion for the total number of successfully parsed skills
	expectedSkillCount := 12
	if parsedCount != expectedSkillCount {
		t.Errorf("Expected to parse %d skills, but got %d", expectedSkillCount, parsedCount)
	}

	// Specific assertions for a known skill
	if artifactsBuilderSkill == nil {
		t.Fatal("Expected to parse 'artifacts-builder' skill, but it was not found or failed parsing.")
	}

	if artifactsBuilderSkill.Meta.Name != "artifacts-builder" {
		t.Errorf("Expected skill name to be 'artifacts-builder', but got '%s'", artifactsBuilderSkill.Meta.Name)
	}

	if len(artifactsBuilderSkill.Resources.Scripts) == 0 {
		t.Error("Expected 'artifacts-builder' to have script resources, but none were found.")
	}

	t.Logf("Successfully parsed %d skills and validated 'artifacts-builder' skill.", parsedCount)
}
