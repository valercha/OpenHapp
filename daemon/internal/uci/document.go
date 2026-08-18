package uci

import (
	"fmt"
	"strings"
)

// Section represents one UCI config section while preserving unknown content.
type Section struct {
	Type   string
	Name   string
	Header string
	Lines  []string
}

// Document represents a parsed UCI configuration file.
type Document struct {
	Prefix   []string
	Sections []Section
}

// ParseDocument parses a UCI configuration while preserving unrelated sections
// and unknown lines.
func ParseDocument(data []byte) (Document, error) {
	lines := strings.Split(string(data), "\n")
	doc := Document{}

	var current *Section

	flush := func() {
		if current == nil {
			return
		}
		doc.Sections = append(doc.Sections, *current)
		current = nil
	}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "config ") {
			flush()

			fields := strings.Fields(trimmed)
			if len(fields) < 2 {
				return Document{}, fmt.Errorf("invalid UCI config line: %q", line)
			}

			section := Section{
				Type:   fields[1],
				Header: line,
			}

			if len(fields) >= 3 {
				section.Name = strings.Trim(fields[2], "'\"")
			}

			current = &section
			continue
		}

		if current == nil {
			doc.Prefix = append(doc.Prefix, line)
		} else {
			current.Lines = append(current.Lines, line)
		}
	}

	flush()

	return doc, nil
}

// Render serializes the document back to UCI text.
func (d Document) Render() []byte {
	var content strings.Builder

	for _, line := range d.Prefix {
		content.WriteString(line)
		content.WriteByte('\n')
	}

	for _, section := range d.Sections {
		content.WriteString(section.Header)
		content.WriteByte('\n')

		for _, line := range section.Lines {
			content.WriteString(line)
			content.WriteByte('\n')
		}
	}

	return []byte(content.String())
}

// Option returns the first scalar option value from the section.
func (s Section) Option(key string) (string, bool) {
	for _, line := range s.Lines {
		fields := strings.Fields(strings.TrimSpace(line))
		if len(fields) < 3 || fields[0] != "option" || fields[1] != key {
			continue
		}

		value := strings.TrimSpace(strings.Join(fields[2:], " "))
		value = strings.Trim(value, "'\"")
		return value, true
	}

	return "", false
}

// SetOption replaces an existing option or appends it to the section.
func (s *Section) SetOption(key, value string) {
	line := fmt.Sprintf("\toption %s '%s'", key, escapeUCIValue(value))

	for i, existing := range s.Lines {
		fields := strings.Fields(strings.TrimSpace(existing))
		if len(fields) >= 2 && fields[0] == "option" && fields[1] == key {
			s.Lines[i] = line
			return
		}
	}

	s.Lines = append(s.Lines, line)
}

// DeleteOption removes all occurrences of an option.
func (s *Section) DeleteOption(key string) {
	filtered := s.Lines[:0]

	for _, line := range s.Lines {
		fields := strings.Fields(strings.TrimSpace(line))
		if len(fields) >= 2 && fields[0] == "option" && fields[1] == key {
			continue
		}
		filtered = append(filtered, line)
	}

	s.Lines = filtered
}

// NewSection creates a named UCI section.
func NewSection(sectionType, name string) Section {
	return Section{
		Type:   sectionType,
		Name:   name,
		Header: fmt.Sprintf("config %s '%s'", sectionType, escapeUCIValue(name)),
	}
}
