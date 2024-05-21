import React, { useState, useEffect, useRef } from "react"
import { Button, ButtonClickEvent, Input, Form, TextArea } from "@fider/components"
import { SignInModal } from "@fider/components"
import { cache, actions, Failure } from "@fider/services"
// import { ImageUpload } from "@fider/models"
import { useFider } from "@fider/hooks"
import { t, Trans } from "@lingui/macro"

interface BoardInputProps {
  placeholder: string
  onTitleChanged: (title: string) => void
}

const CACHE_TITLE_KEY = "PostInput-Title"
const CACHE_DESCRIPTION_KEY = "PostInput-Description"

export const BoardInput = (props: BoardInputProps) => {
  const getCachedValue = (key: string): string => {
    if (fider.session.isAuthenticated) {
      return cache.session.get(key) || ""
    }
    return ""
  }

  const fider = useFider()
  const titleRef = useRef<HTMLInputElement>()
  const [title, setTitle] = useState(getCachedValue(CACHE_TITLE_KEY))
  const [description, setDescription] = useState(getCachedValue(CACHE_DESCRIPTION_KEY))
  const [isSignInModalOpen, setIsSignInModalOpen] = useState(false)
  // const [attachments, setAttachments] = useState<ImageUpload[]>([])
  const [error, setError] = useState<Failure | undefined>(undefined)

  useEffect(() => {
    props.onTitleChanged(title)
  }, [title])

  const handleTitleFocus = () => {
    if (!fider.session.isAuthenticated && titleRef.current) {
      titleRef.current.blur()
      setIsSignInModalOpen(true)
    }
  }

  const handleTitleChange = (value: string) => {
    cache.session.set(CACHE_TITLE_KEY, value)
    setTitle(value)
    props.onTitleChanged(value)
  }

  const hideModal = () => setIsSignInModalOpen(false)
  const clearError = () => setError(undefined)

  const handleDescriptionChange = (value: string) => {
    cache.session.set(CACHE_DESCRIPTION_KEY, value)
    setDescription(value)
  }

  const submit = async (event: ButtonClickEvent) => {
    if (title) {
      const result = await actions.createBoard(title, description)
      if (result.ok) {
        clearError()
        console.log("Board created")
        console.log(result)
        console.log(result.data)
        cache.session.remove(CACHE_TITLE_KEY, CACHE_DESCRIPTION_KEY)
        location.href = `board/${result.data.id}`
        event.preventEnable()
      } else if (result.error) {
        setError(result.error)
      }
    }
  }

  const details = () => (
    <>
      <TextArea
        field="welcome"
        onChange={handleDescriptionChange}
        value={description}
        minRows={5}
        placeholder={t({ id: "home.boardinput.description.placeholder", message: "Provide a welcome message (optional)" })}
      />
      {/* <MultiImageUploader field="attachments" maxUploads={3} onChange={setAttachments} /> */}
      <Button type="submit" variant="primary" onClick={submit}>
        <Trans id="action.submit">Submit</Trans>
      </Button>
    </>
  )

  return (
    <>
      <SignInModal isOpen={isSignInModalOpen} onClose={hideModal} />
      <Form error={error}>
        <Input
          field="title"
          disabled={fider.isReadOnly}
          noTabFocus={!fider.session.isAuthenticated}
          inputRef={titleRef}
          onFocus={handleTitleFocus}
          maxLength={100}
          value={title}
          onChange={handleTitleChange}
          placeholder={props.placeholder}
        />
        {title && details()}
      </Form>
    </>
  )
}
